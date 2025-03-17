# set domain name here
variable "domain_name" {
  type    = string
  default = "tadeutupinamba.com.br"
}

resource "aws_route53_zone" "primary" {
  name = var.domain_name
}

# Create Route53 record for the API endpoint
resource "aws_route53_record" "api" {
  zone_id = aws_route53_zone.primary.zone_id # Reference existing zone resource
  name    = "api.${var.domain_name}"
  type    = "A"

  alias {
    name                   = kubernetes_ingress_v1.fastfood_ingress.status.0.load_balancer.0.ingress.0.hostname
    zone_id                = "Z35SXDOTRQ7X7K" # This is the fixed zone ID for us-east-1 ALBs
    evaluate_target_health = true
  }
}

resource "aws_acm_certificate" "api_cert" {
  domain_name       = "api.${var.domain_name}"
  validation_method = "DNS"
  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

resource "aws_route53_record" "api_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.api_cert.domain_validation_options : dvo.domain_name => {
      name  = dvo.resource_record_name
      type  = dvo.resource_record_type
      value = dvo.resource_record_value
    }
  }

  zone_id = aws_route53_zone.primary.zone_id
  name    = each.value.name
  type    = each.value.type
  records = [each.value.value]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "api_cert_validation" {
  certificate_arn         = aws_acm_certificate.api_cert.arn
  validation_record_fqdns = [for record in aws_route53_record.api_cert_validation : record.fqdn]
}

# AWS Load Balancer Controller IAM Policy
resource "aws_iam_policy" "aws_load_balancer_controller" {
  name        = "AWSLoadBalancerControllerIAMPolicy"
  description = "Policy for AWS Load Balancer Controller"

  policy = file("${path.module}/iam-policy.json")
  # You'll need to download this file:
  # curl -o iam-policy.json https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/main/docs/install/iam_policy.json
}

# IAM Role for Service Account (IRSA)
module "lb_controller_role" {
  source                        = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
  version                       = "~> 4.0"
  create_role                   = true
  role_name                     = "eks-alb-controller"
  provider_url                  = replace(module.eks.cluster_oidc_issuer_url, "https://", "")
  role_policy_arns              = [aws_iam_policy.aws_load_balancer_controller.arn]
  oidc_fully_qualified_subjects = ["system:serviceaccount:kube-system:aws-load-balancer-controller"]
}

# Kubernetes service account for the controller
resource "kubernetes_service_account" "aws_load_balancer_controller" {
  metadata {
    name      = "aws-load-balancer-controller"
    namespace = "kube-system"
    labels = {
      "app.kubernetes.io/name"      = "aws-load-balancer-controller"
      "app.kubernetes.io/component" = "controller"
    }
    annotations = {
      "eks.amazonaws.com/role-arn" = module.lb_controller_role.iam_role_arn
    }
  }
}
resource "helm_release" "aws_load_balancer_controller" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  namespace  = "kube-system"

  set {
    name  = "clusterName"
    value = var.cluster_name
  }

  set {
    name  = "serviceAccount.create"
    value = "false"
  }

  set {
    name  = "serviceAccount.name"
    value = "aws-load-balancer-controller"
  }

  set {
    name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
    value = module.lb_controller_role.iam_role_arn
  }

  depends_on = [
    module.eks,
    kubernetes_service_account.aws_load_balancer_controller
  ]
}

# IngressClass resource
# resource "kubernetes_ingress_class_v1" "alb" {
#   depends_on = [helm_release.aws_load_balancer_controller]

#   metadata {
#     name = "alb"
#   }

#   spec {
#     controller = "ingress.k8s.aws/alb"
#   }
# }

resource "kubernetes_ingress_v1" "fastfood_ingress" {
  depends_on = [helm_release.aws_load_balancer_controller]

  metadata {
    name = "fastfood-ingress"
    annotations = {
      "alb.ingress.kubernetes.io/scheme"          = "internet-facing"
      "alb.ingress.kubernetes.io/target-type"     = "ip"
      "alb.ingress.kubernetes.io/listen-ports"    = "[{\"HTTPS\":443}, {\"HTTP\":80}]"
      "alb.ingress.kubernetes.io/certificate-arn" = aws_acm_certificate.api_cert.arn
      "alb.ingress.kubernetes.io/ssl-redirect"    = "443"
    }
  }

  spec {
    ingress_class_name = "alb"

    rule {
      host = "api.${var.domain_name}"

      http {
        path {
          path      = "/"
          path_type = "Prefix"

          backend {
            service {
              name = "restaurant-api-service"

              port {
                number = 80
              }
            }
          }
        }
      }
    }
  }
}
