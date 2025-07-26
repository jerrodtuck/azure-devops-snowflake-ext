import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as SDK from 'azure-devops-extension-sdk';
import { SnowflakeDropdown } from './components/SnowflakeDropdown/SnowflakeDropdown';
import './index.css';

// Initialize the Azure DevOps SDK
SDK.init({
  loaded: false,
  applyTheme: true
});

// Wait for SDK to be ready
SDK.ready().then(() => {
  // Register the control with the full contribution ID
  SDK.register(SDK.getContributionId(), () => {
    const config = SDK.getConfiguration();
    
    return {
      // Called when the control is first loaded
      onLoaded: async () => {
        // Get the container element
        const container = document.getElementById('root');
        if (!container) return;

        // Get the work item form service
        const workItemFormService = await SDK.getService<any>(
          "ms.vss-work-web.work-item-form"
        );

        // Get configuration from Azure DevOps
        const fieldName = config.witInputs?.FieldName || '';
        const apiUrl = config.witInputs?.ApiUrl || 'http://localhost:8080/api';
        const dataType = config.witInputs?.DataType || 'cc';
        const minSearchLength = config.witInputs?.MinSearchLength || 2;
        const debounceDelay = config.witInputs?.DebounceDelay || 300;

        // Get the current field value
        let currentValue = '';
        try {
          currentValue = await workItemFormService.getFieldValue(fieldName) || '';
        } catch (e) {
          console.error('Error getting field value:', e);
        }

        // Render the dropdown component
        ReactDOM.render(
          <SnowflakeDropdown
            fieldName={fieldName}
            dataType={dataType}
            apiUrl={apiUrl}
            minSearchLength={minSearchLength}
            debounceDelay={debounceDelay}
            initialValue={currentValue}
            onValueChange={async (value) => {
              // Notify Azure DevOps of value change
              try {
                await workItemFormService.setFieldValue(fieldName, value);
              } catch (e) {
                console.error('Error setting field value:', e);
              }
            }}
          />,
          container
        );
      },

      // Called when the field value changes
      onFieldChanged: (args: any) => {
        if (args && args.changedFields && args.changedFields[config.witInputs?.FieldName]) {
          // Field value changed externally
          // You could update the component here if needed
        }
      },

      // Get the current value
      getValue: () => {
        // This would need to be implemented with a ref or state management
        return '';
      },

      // Called when the control is being disposed
      onUnloaded: () => {
        const container = document.getElementById('root');
        if (container) {
          ReactDOM.unmountComponentAtNode(container);
        }
      }
    };
  });
  
  // Notify that the extension is loaded
  SDK.notifyLoadSucceeded();
});
