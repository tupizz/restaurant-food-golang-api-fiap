# variable "route53_zone_id" {
#   description = "The Route 53 hosted zone ID for your domain"
#   type        = string
#   default     = "Z35SXDOTRQ7X7K"
# }
provider "aws" {
  alias  = "acm_region"
  region = var.aws_region # Ensure this matches your API Gateway region
}


variable "domain_name" {
  description = "The domain name for the API Gateway"
  type        = string
  default     = "tadeutupinamba.com.br"
}

# Look up the hosted zone ID for the domain
data "aws_route53_zone" "selected" {
  name         = var.domain_name
  private_zone = false
}


# ACM Certificate for the API Gateway domain
resource "aws_acm_certificate" "api_cert" {
  provider          = aws.acm_region
  domain_name       = "auth.${var.domain_name}"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = local.tags
}

# DNS validation records for the certificate
resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.api_cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.selected.zone_id
}

# Certificate validation
resource "aws_acm_certificate_validation" "api_cert" {
  provider                = aws.acm_region
  certificate_arn         = aws_acm_certificate.api_cert.arn
  validation_record_fqdns = [for record in aws_route53_record.cert_validation : record.fqdn]
}

# Custom domain for API Gateway
resource "aws_apigatewayv2_domain_name" "api_domain" {
  domain_name = "auth.${var.domain_name}"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.api_cert.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  tags = local.tags
}

# API mapping to connect the custom domain to your API stage
resource "aws_apigatewayv2_api_mapping" "api_mapping" {
  api_id      = aws_apigatewayv2_api.lambda_api.id
  domain_name = aws_apigatewayv2_domain_name.api_domain.id
  stage       = aws_apigatewayv2_stage.lambda_stage.id
}

# Route 53 A record pointing to the API Gateway
resource "aws_route53_record" "api_record" {
  name    = "auth.${var.domain_name}"
  type    = "A"
  zone_id = data.aws_route53_zone.selected.zone_id

  alias {
    name                   = aws_apigatewayv2_domain_name.api_domain.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.api_domain.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

# Output the custom domain URL
output "api_custom_domain" {
  value = "https://${aws_apigatewayv2_domain_name.api_domain.domain_name}"
}
