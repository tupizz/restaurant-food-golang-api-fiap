   # Destroy ingress first
   terraform destroy -target=kubernetes_ingress_v1.fastfood_ingress -auto-approve
   
   # Then Helm release
   terraform destroy -target=helm_release.aws_load_balancer_controller -auto-approve
   
   # Then service account
   terraform destroy -target=kubernetes_service_account.aws_load_balancer_controller -auto-approve

   terraform destroy -target=aws_route53_record.api -auto-approve
   terraform destroy -target=aws_acm_certificate_validation.api_cert_validation -auto-approve
   terraform destroy -target=aws_route53_record.api_cert_validation -auto-approve
   terraform destroy -target=aws_acm_certificate.api_cert -auto-approve
      # First addons
   terraform destroy -target=aws_eks_addon.addons -auto-approve
   terraform destroy -target=aws_eks_addon.ebs_csi_driver -auto-approve
   
   # Then the EKS cluster
   terraform destroy -target=module.eks -auto-approve

    # Destroy VPC
   terraform destroy -target=module.vpc -auto-approve
   
   # Finally, destroy everything else
   terraform destroy -auto-approve