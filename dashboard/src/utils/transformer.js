export function transformResourceAuditPayload(inputPayload) {
  try {
    if (!inputPayload || !inputPayload['0'] || !inputPayload['0'].credentials) {
      throw new Error('Invalid input payload structure');
    }
    let parsedCredentials;
    try {
      parsedCredentials = JSON.parse(inputPayload['0'].credentials);
    } catch (error) {
      throw new Error('Failed to parse credentials JSON: ' + error.message);
    }
    return {
      name: inputPayload['0'].name,
      provider: 'gcp',
      credentials: parsedCredentials,
    };
  } catch (error) {
    console.error('Transformation error:', error);
    return { error: error.message };
  }
}

export function transformResourceAuditData(inputData) {
  return {
    data: inputData.map((account) => {
      // Extract date and format it
      const updatedAt = new Date(account.updatedAt);
      const formattedDate =
        updatedAt.toLocaleDateString('en-US', {
          day: 'numeric',
          month: 'long',
          year: 'numeric',
        }) +
        ', ' +
        updatedAt.toLocaleTimeString('en-US', {
          hour: '2-digit',
          minute: '2-digit',
          hour12: false,
        });

      // Initialize audit data structure
      const initialCounts = {
        danger: 0,
        warning: 0,
        pending: 0,
        compliant: 0,
        unchecked: 0,
        total: 0,
      };

      // Default auditData
      let auditData = {};

      // Process audit details if available and not empty
      if (account?.auditDetails && account?.auditDetails?.data) {
        const auditDetailsData = account.auditDetails?.data;

        // Check if auditDetailsData is not empty
        if (Object.keys(auditDetailsData).length > 0) {
          // Initialize the 'all' category
          auditData.all = { ...initialCounts };

          if (Array.isArray(auditDetailsData)) {
            // Process array format - we don't have category information
            auditDetailsData.forEach((audit) => {
              const items = audit?.result?.items || [];

              items.forEach((item) => {
                // Map status to standardized format
                let status;
                if (item.status === 'failing' || item.status === 'danger') {
                  status = 'danger';
                } else if (item.status === 'warning') {
                  status = 'warning';
                } else if (item.status === 'pending') {
                  status = 'pending';
                } else if (item.status === 'compliant' || item.status === 'passing') {
                  status = 'compliant';
                } else {
                  status = 'unchecked';
                }

                // Increment the 'all' category counter
                auditData.all[status]++;
                auditData.all.total++;
              });
            });
          } else {
            // Process object format with category keys - these are the actual categories
            Object.keys(auditDetailsData).forEach((category) => {
              // Initialize this category
              if (!auditData[category]) {
                auditData[category] = { ...initialCounts };
              }

              const categoryData = auditDetailsData[category];

              if (Array.isArray(categoryData)) {
                categoryData.forEach((audit) => {
                  const items = audit.result.items || [];

                  items.forEach((item) => {
                    // Map status to standardized format
                    let status;
                    if (item.status === 'failing' || item.status === 'danger') {
                      status = 'danger';
                    } else if (item.status === 'warning') {
                      status = 'warning';
                    } else if (item.status === 'pending') {
                      status = 'pending';
                    } else if (item.status === 'compliant' || item.status === 'passing') {
                      status = 'compliant';
                    } else {
                      status = 'unchecked';
                    }

                    // Increment the category counter
                    auditData[category][status]++;
                    auditData[category].total++;

                    // Also increment the 'all' category
                    auditData.all[status]++;
                    auditData.all.total++;
                  });
                });
              }
            });
          }
        }
        // If auditDetailsData is empty, auditData will remain an empty object
      }

      const transformedAccount = {
        ...account,
        id: account.id.toString(),
        status: 'READY',
        icon: 'cloud',
        updatedAt: formattedDate,
        auditData: auditData,
        categoryIcons: {
          stale: 'server',
          overprovision: 'exclamation',
          security: 'shield',
        },
      };

      return transformedAccount;
    }),
  };
}
