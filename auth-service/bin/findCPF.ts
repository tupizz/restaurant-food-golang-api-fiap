#!/usr/bin/env ts-node

import {
  AdminGetUserCommand,
  CognitoIdentityProviderClient,
} from "@aws-sdk/client-cognito-identity-provider";

// Configuration
const REGION = process.env.AWS_REGION || "us-east-1";
const USER_POOL_ID = "us-east-1_Sy6UT9SLi";
const CPF = process.argv[2] || "52998224725"; // Default CPF or take from command line argument

// Initialize the Cognito client
const cognitoClient = new CognitoIdentityProviderClient({
  region: REGION,
});

// Function to find user by CPF (now directly using AdminGetUser)
const findUserByCPF = async (cpf: string) => {
  try {
    console.log(`Looking up user with CPF (username): ${cpf}`);

    try {
      const command = new AdminGetUserCommand({
        UserPoolId: USER_POOL_ID,
        Username: cpf,
      });

      const user = await cognitoClient.send(command);

      console.log("\nUser found!");
      console.log(`Username: ${user.Username}`);
      console.log("Attributes:");
      user.UserAttributes?.forEach((attr) => {
        console.log(`  ${attr.Name}: ${attr.Value}`);
      });

      return user;
    } catch (error: any) {
      if (error.name === "UserNotFoundException") {
        console.log("User not found with this CPF");
        return null;
      }
      throw error;
    }
  } catch (error) {
    console.error("Error finding user by CPF:", error);
    throw error;
  }
};

// Main function
const main = async () => {
  try {
    await findUserByCPF(CPF);
  } catch (error) {
    console.error("Script failed:", error);
    process.exit(1);
  }
};

// Run the script
main();
