import React from 'react';

const headlineStyle = 'text-gray-600 font-medium text-base';

export default function AwsAccountCreationGuide() {
  return (
    <div className="w-full flex flex-col space-y-4">
      <span className={headlineStyle}>1. Sign In to AWS Management Console</span>
      <span className="text-gray-500 pl-4 list-disc">
        • Log in to the AWS Management Console using your AWS account credentials. Ensure you have
        the necessary permissions to create and manage resources.
      </span>

      <span className={`${headlineStyle} my-2`}>2. Open the IAM Service</span>
      <span className="text-gray-500 pl-4 list-disc">
        • In the AWS Management Console, search for IAM (Identity and Access Management) in the
        search bar and open the IAM dashboard.
      </span>

      <span className={`${headlineStyle} my-2`}>3. Navigate to Users Section</span>
      <span className="text-gray-500 pl-4 list-disc">
        • On the IAM dashboard, click on Users in the left-hand navigation pane to view and manage
        users in your AWS account.
      </span>

      <span className={`${headlineStyle} my-2`}>4. Add a New User</span>
      <span className="text-gray-500 pl-4 list-disc">
        • Click on the Add user button to start creating a new IAM user.
      </span>

      <span className={`${headlineStyle} my-2`}>5. Provide User Details</span>
      <span className="text-gray-500 pl-4 list-disc">
        • Enter a unique name in the User name field. This name should clearly identify the purpose
        or owner of the account (e.g., Zop-Admin).
      </span>

      <span className={`${headlineStyle} my-2`}>6. Set Permissions</span>
      <span className="text-gray-500 pl-4 list-disc">
        • Create a User Group with the{' '}
        <span className="font-bold">Administrator Access policy</span> . Attach this group to the
        user for full access to AWS resources.
      </span>

      <span className={`${headlineStyle} my-2`}>7. Complete the User Creation Process</span>
      <span className="text-gray-500 pl-4 list-disc">
        • Click on the Create user button to finalize the creation of the IAM user. The user will
        now appear in your IAM user list.
      </span>

      <span className={`${headlineStyle} my-2`}>8. Generate and Save User Credentials</span>
      <span className="text-gray-500 pl-4 list-disc">
        • After creating the user, you’ll be prompted to download the Access Key ID and Secret
        Access Key. Save these credentials securely, and use them for connecting your aws account
        with zop.dev.
      </span>
    </div>
  );
}
