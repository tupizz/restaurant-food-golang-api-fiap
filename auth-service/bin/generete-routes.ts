import * as fs from "fs";
import { generateTerraformRouteKeys } from "../src/routes";

// Generate Terraform configuration for routes
function generateTerraformRoutes(): string {
  const routeKeys = generateTerraformRouteKeys();

  let terraformConfig = `# This file is generated automatically. Do not edit directly.\n\n`;

  routeKeys.forEach((routeKey, index) => {
    terraformConfig += `resource "aws_apigatewayv2_route" "route_${index}" {\n`;
    terraformConfig += `  api_id    = aws_apigatewayv2_api.lambda_api.id\n`;
    terraformConfig += `  route_key = "${routeKey}"\n`;
    terraformConfig += `  target    = "integrations/\${aws_apigatewayv2_integration.lambda_integration.id}"\n`;
    terraformConfig += `}\n\n`;
  });

  return terraformConfig;
}

// Write the generated routes to a file
const terraformRoutes = generateTerraformRoutes();
fs.writeFileSync("terraform/routes.tf", terraformRoutes);
console.log("Routes generated successfully in terraform/routes.tf");
