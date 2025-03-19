# Define AWS as the provider for this infrastructure
# This tells Terraform to use AWS APIs for creating resources
provider "aws" {
  region = var.aws_region # Use the region defined in variables
}

# Define required Terraform providers and versions
# This ensures compatibility and consistent behavior
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws" # Official AWS provider from HashiCorp
      version = "~> 5.0"        # Use version 5.x
    }
    kubernetes = {
      source  = "hashicorp/kubernetes" # Provider to manage Kubernetes resources
      version = "~> 2.10"
    }
    helm = {
      source  = "hashicorp/helm" # Provider to deploy Helm charts in Kubernetes
      version = "~> 2.5"
    }
  }

  # Configure Terraform backend to store state in S3
  # This allows team collaboration and state persistence
  # Without this, Terraform state would be stored locally, making team collaboration difficult
  backend "s3" {
    bucket = "fiap-tf-state-bucket"       # S3 bucket to store Terraform state
    key    = "k8s/terraform-base.tfstate" # Path within the bucket
    region = "us-east-1"                  # Region where the S3 bucket is located
  }
}

# Configure Kubernetes provider to interact with the EKS cluster
# This provider allows Terraform to manage Kubernetes resources
provider "kubernetes" {
  host                   = module.eks.cluster_endpoint                                 # API endpoint of the cluster
  cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data) # Cluster CA cert

  # Authentication method using AWS CLI to get an authentication token
  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "aws"
    args        = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
  }
}

# Configure Helm provider to deploy applications to Kubernetes
# Helm is a package manager for Kubernetes that simplifies application deployment
provider "helm" {
  kubernetes {
    host                   = module.eks.cluster_endpoint
    cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)

    # Same authentication method as the Kubernetes provider
    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      args        = ["eks", "get-token", "--cluster-name", module.eks.cluster_name]
    }
  }
}

# Create VPC (Virtual Private Cloud) for EKS
# A VPC is a virtual network dedicated to your AWS account
# EKS requires a VPC with specific configurations to function properly
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws" # Using official AWS VPC module
  version = "5.0.0"

  name = "${var.project_name}-vpc" # Name the VPC based on project
  cidr = var.vpc_cidr              # IP address range for the VPC

  azs             = var.availability_zones # Availability zones for redundancy
  private_subnets = var.private_subnets    # Private subnets for EKS nodes
  public_subnets  = var.public_subnets     # Public subnets for load balancers

  enable_nat_gateway   = false # NAT gateway allows private subnet resources to access internet
  single_nat_gateway   = true  # Use single NAT to reduce costs (less redundant but cheaper)
  enable_dns_hostnames = true  # Enable DNS hostnames for the VPC
  enable_dns_support   = true  # Enable DNS resolution in the VPC

  # Add this line to enable auto-assign public IPs in public subnets
  map_public_ip_on_launch = true

  # Tags are key-value pairs attached to AWS resources for identification and automation
  tags = {
    Environment                                 = var.environment
    Project                                     = var.project_name
    "kubernetes.io/cluster/${var.cluster_name}" = "shared" # Required tag for EKS to identify VPC resources
  }

  # Special tags for subnets used by EKS
  # These tags are required for EKS to discover and use the subnets correctly
  private_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/internal-elb"           = "1" # Marks subnets for internal load balancers
  }

  public_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                    = "1" # Marks subnets for external load balancers
  }
}

# Create EKS Cluster
# EKS is AWS's managed Kubernetes service
module "eks" {
  source  = "terraform-aws-modules/eks/aws" # Using official AWS EKS module
  version = "~> 19.0"

  cluster_name    = var.cluster_name    # Name of the EKS cluster
  cluster_version = var.cluster_version # Kubernetes version to use

  vpc_id     = module.vpc.vpc_id         # Use the VPC we created above
  subnet_ids = module.vpc.public_subnets # Use public subnets for node groups when NAT gateway is disabled

  cluster_endpoint_public_access = true # Allow public access to cluster API endpoint

  # Define the worker node groups for the cluster
  # These are the EC2 instances that will run your containerized applications
  eks_managed_node_groups = {
    general = {
      desired_size = 2 # CHANGE: Increased from 1 to 2
      min_size     = 1
      max_size     = 3 # CHANGE: Increased from 1 to 3

      instance_types = ["t3.medium"] # CHANGE: Increased from t3.small to t3.medium
      capacity_type  = "SPOT"        # Keeping SPOT for cost savings

      labels = {
        Environment = var.environment
        Project     = var.project_name
      }

      # Attach IAM policies to node group for required AWS permissions
      # These policies allow worker nodes to access various AWS services
      iam_role_additional_policies = {
        AmazonEBSCSIDriverPolicy           = "arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy" # For persistent storage
        AmazonEKSWorkerNodePolicy          = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"             # Core EKS permissions
        AmazonEKS_CNI_Policy               = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"                  # For pod networking
        AmazonEC2ContainerRegistryReadOnly = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"    # To pull container images
        CloudWatchAgentServerPolicy        = "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"           # For monitoring
        AmazonSSMManagedInstanceCore       = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"          # For SSM access
      }
    }
  }

  # IAM configuration for the EKS cluster
  create_iam_role = true                               # Create IAM role for cluster
  iam_role_name   = "${var.cluster_name}-cluster-role" # Name of the IAM role

  # Encryption configuration
  attach_cluster_encryption_policy = true # Encrypt cluster data
  cluster_encryption_policy_name   = "${var.cluster_name}-encryption-policy"

  enable_irsa = true # Enable IAM roles for service accounts
  # This allows Kubernetes pods to have specific AWS permissions

  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

# Attach additional IAM policies to the EKS cluster role
# These provide the cluster control plane with necessary AWS permissions
resource "aws_iam_role_policy_attachment" "eks_cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy" # Core permissions for EKS management
  role       = module.eks.cluster_iam_role_name
}

resource "aws_iam_role_policy_attachment" "eks_vpc_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController" # For managing VPC resources
  role       = module.eks.cluster_iam_role_name
}

# Retry settings for cluster operations
# Ensures more reliable execution of operations
locals {
  retry_join_enabled = true
  max_retry_attempts = 3 # Maximum number of times to retry connecting to the cluster
}

# Check that the cluster is available after creation
# This helps ensure that subsequent resources depending on the cluster can be created
resource "null_resource" "cluster_check" {
  depends_on = [module.eks] # Wait for EKS cluster to be created first

  provisioner "local-exec" {
    command = <<-EOT
      for i in $(seq 1 ${local.max_retry_attempts}); do
        if aws eks describe-cluster --name ${var.cluster_name} --region ${var.aws_region}; then
          exit 0
        fi
        sleep 30
      done
      exit 1
    EOT
  }
}

# Get AWS account ID for use in OIDC provider configuration
data "aws_caller_identity" "current" {}

# Install EKS add-ons
# Add-ons provide essential functionality to the Kubernetes cluster
resource "aws_eks_addon" "addons" {
  for_each = {
    coredns    = {} # DNS service for Kubernetes service discovery
    kube-proxy = {} # Network proxy for Kubernetes service abstraction
  }

  cluster_name = module.eks.cluster_name
  addon_name   = each.key

  # Conflict resolution settings
  # OVERWRITE on create means we'll replace any existing versions
  # PRESERVE on update means we won't override custom configurations during updates
  resolve_conflicts_on_create = "OVERWRITE"
  resolve_conflicts_on_update = "PRESERVE"

  depends_on = [module.eks]
}

# Add EBS CSI driver separately (requires special IAM configuration)
# This driver allows Kubernetes pods to use EBS volumes for persistent storage
resource "aws_eks_addon" "ebs_csi_driver" {
  cluster_name             = module.eks.cluster_name
  addon_name               = "aws-ebs-csi-driver"
  service_account_role_arn = aws_iam_role.ebs_csi_driver.arn # Use the IAM role defined below

  # Same conflict resolution settings as other add-ons
  resolve_conflicts_on_create = "OVERWRITE"
  resolve_conflicts_on_update = "PRESERVE"

  depends_on = [module.eks]
}

# Create IAM role for AWS Load Balancer Controller
# This controller manages AWS load balancers for Kubernetes services
resource "aws_iam_role" "lb_controller" {
  name = "${var.cluster_name}-lb-controller"

  # Trust policy allowing the Kubernetes service account to assume this role
  # This uses OIDC federation, which links Kubernetes service accounts with AWS IAM
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/${module.eks.oidc_provider}"
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "${module.eks.oidc_provider}:sub" = "system:serviceaccount:kube-system:aws-load-balancer-controller"
          }
        }
      }
    ]
  })
}

# Policy for the Load Balancer Controller
# These permissions allow managing AWS load balancers from Kubernetes
resource "aws_iam_policy" "lb_controller" {
  name        = "${var.cluster_name}-lb-controller-policy"
  description = "Policy for AWS Load Balancer Controller"

  # Extensive policy allowing the controller to manage load balancers, target groups, etc.
  policy = aws_iam_policy.aws_load_balancer_controller.policy
}

# Create IAM role for Cluster Autoscaler
# This component automatically adjusts the size of node groups based on demand
resource "aws_iam_role" "cluster_autoscaler" {
  name = "${var.cluster_name}-cluster-autoscaler"

  # Similar OIDC trust policy but for the autoscaler service account
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/${module.eks.oidc_provider}"
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "${module.eks.oidc_provider}:sub" = "system:serviceaccount:kube-system:cluster-autoscaler"
          }
        }
      }
    ]
  })
}

# Policy for Cluster Autoscaler
# These permissions allow scaling node groups up and down
resource "aws_iam_policy" "cluster_autoscaler" {
  name        = "${var.cluster_name}-cluster-autoscaler-policy"
  description = "Policy for Cluster Autoscaler"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "autoscaling:DescribeAutoScalingGroups",
        "autoscaling:DescribeAutoScalingInstances",
        "autoscaling:DescribeLaunchConfigurations",
        "autoscaling:DescribeScalingActivities",
        "autoscaling:DescribeTags",
        "ec2:DescribeInstanceTypes",
        "ec2:DescribeLaunchTemplateVersions"
      ],
      "Resource": ["*"]
    },
    {
      "Effect": "Allow",
      "Action": [
        "autoscaling:SetDesiredCapacity",
        "autoscaling:TerminateInstanceInAutoScalingGroup",
        "ec2:DescribeImages",
        "ec2:GetInstanceTypesFromInstanceRequirements",
        "eks:DescribeNodegroup"
      ],
      "Resource": ["*"]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "cluster_autoscaler" {
  policy_arn = aws_iam_policy.cluster_autoscaler.arn
  role       = aws_iam_role.cluster_autoscaler.name
}

# Create IAM role for External DNS
# This component automatically manages DNS records for Kubernetes services
resource "aws_iam_role" "external_dns" {
  name = "${var.cluster_name}-external-dns"

  # OIDC trust policy for External DNS service account
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/${module.eks.oidc_provider}"
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "${module.eks.oidc_provider}:sub" = "system:serviceaccount:kube-system:external-dns"
          }
        }
      }
    ]
  })
}

# Policy for External DNS
# These permissions allow managing Route 53 DNS records
resource "aws_iam_policy" "external_dns" {
  name        = "${var.cluster_name}-external-dns-policy"
  description = "Policy for External DNS"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "route53:ChangeResourceRecordSets"
      ],
      "Resource": [
        "arn:aws:route53:::hostedzone/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "route53:ListHostedZones",
        "route53:ListResourceRecordSets"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "external_dns" {
  policy_arn = aws_iam_policy.external_dns.arn
  role       = aws_iam_role.external_dns.name
}

# Create IAM role for accessing AWS Secrets Manager
# This allows securely storing and retrieving sensitive configuration
resource "aws_iam_role" "secrets_manager" {
  name = "${var.cluster_name}-secrets-manager"

  # OIDC trust policy for External Secrets service account
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/${module.eks.oidc_provider}"
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "${module.eks.oidc_provider}:sub" = "system:serviceaccount:kube-system:external-secrets"
          }
        }
      }
    ]
  })
}

# Policy for Secrets Manager access
# These permissions allow reading secrets from AWS Secrets Manager
resource "aws_iam_policy" "secrets_manager" {
  name        = "${var.cluster_name}-secrets-manager-policy"
  description = "Policy for Secrets Manager access"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetResourcePolicy",
        "secretsmanager:GetSecretValue",
        "secretsmanager:DescribeSecret",
        "secretsmanager:ListSecretVersionIds"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "secrets_manager" {
  policy_arn = aws_iam_policy.secrets_manager.arn
  role       = aws_iam_role.secrets_manager.name
}

# Create IAM role for EBS CSI Driver
# This driver allows Kubernetes pods to use EBS volumes for persistent storage
resource "aws_iam_role" "ebs_csi_driver" {
  name = "${var.cluster_name}-ebs-csi-driver"

  # OIDC trust policy for EBS CSI Driver service account
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/${module.eks.oidc_provider}"
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "${module.eks.oidc_provider}:sub" = "system:serviceaccount:kube-system:ebs-csi-controller-sa"
          }
        }
      }
    ]
  })
}

# Attach managed EBS CSI Driver policy to the role
# This AWS-managed policy provides the necessary permissions to manage EBS volumes
resource "aws_iam_role_policy_attachment" "ebs_csi_driver" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"
  role       = aws_iam_role.ebs_csi_driver.name
}
