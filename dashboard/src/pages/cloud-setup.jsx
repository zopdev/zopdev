import React, { useEffect } from 'react';
import CreateCloud from '@/components/organisms/CreateCloud.jsx';
import DynamicFormRadioWithIcon from '@/components/atom/Button/RadioButtonWithIcon/index.jsx';
import Stepper from '@/components/organisms/Stepper.jsx';
import { usePostAuditData } from '@/queries/CloudAccount/index.js';
import { transformResourceAuditPayload } from '@/utils/transformer.js';

const ResourceAudit = ({ data, updateData, setIsComplete }) => {
  const auditOptions = [
    {
      label: 'Stale',
      description: 'Identify the resources that are no longer in use.',
      value: 'Stale',
    },
    {
      label: 'Overprovision',
      description: 'Resources that have more capacity than needed',
      value: 'Overprovision',
    },
    {
      label: 'Security',
      description: 'Resources with potential security issues',
      value: 'Security',
    },
    {
      label: 'Run All',
      description: 'Run all types of audits on your resources',
      value: 'run-all',
    },
  ];
  const handleChange = (newValue) => {
    const updatedData = { ...data, selectedOption: newValue };
    updateData(updatedData);
    validateStep(updatedData);
  };

  const validateStep = (stepData) => {
    const isValid = !!stepData?.selectedOption;
    setIsComplete(isValid);
  };

  useEffect(() => {
    if (!data?.selectedOption) {
      updateData({ ...data, selectedOption: 'run-all' });
    } else {
      validateStep(data);
    }
  }, [data, updateData, setIsComplete]);

  return (
    <div className="md:flex xs:space-y-12 space-x-8 min-h-[28rem] mt-10 ml-2">
      <div className="md:w-[34%] md:m-12 md:mx-14">
        <h2 className="text-md font-semibold mb-2">Select an Audit Type</h2>
        <p className="text-secondary-600">
          Evaluate your cloud environment for security risks, performance bottlenecks, or cost
          inefficiencies.
        </p>
      </div>
      <div className="md:w-[45%] md:my-14">
        <DynamicFormRadioWithIcon
          options={auditOptions}
          name="resource-audit"
          value={data?.selectedOption}
          defaultSelected={data?.selectedOption || 'run-all'}
          onChange={handleChange}
          orientation="horizontal"
        />
      </div>
    </div>
  );
};

// const ScheduleStep = ({ data, updateData, setIsComplete }) => {
//   const cronPattern =
//     /(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|sun|mon|tue|wed|thu|fri|sat|\*) ?){5,7})/;
//   const [error, setError] = useState('');
//
//   const handleChange = (e) => {
//     const { value } = e.target;
//     const updatedData = { ...data, cronSchedule: value };
//     updateData(updatedData);
//     validateStep(updatedData);
//   };
//
//   const validateStep = (stepData) => {
//     const value = stepData?.cronSchedule || '';
//     if (value.trim() === '') {
//       setError('');
//       setIsComplete(false);
//     } else if (!cronPattern.test(value)) {
//       setError('Invalid Schedule Time');
//       setIsComplete(false);
//     } else {
//       setError('');
//       setIsComplete(true);
//     }
//   };
//
//   useEffect(() => {
//     validateStep(data);
//   }, [data, setIsComplete]);
//
//   return (
//     <div className="md:flex xs:space-y-12 space-x-8 min-h-[28rem] mt-10 ml-2">
//       <div className="md:w-[34%] md:m-12 md:mx-14">
//         <h2 className="text-md font-semibold mb-2">Give Audit Schedule</h2>
//         <p className="text-secondary-600 mb-2">
//           Specify how often you want to run the resource audit by entering a time interval.
//         </p>
//         <p className="text-sm text-secondary-500">
//           You can use special formats like @daily, @hourly, or @every 5m to set the frequency.
//         </p>
//       </div>
//
//       <div className="w-[35%] space-y-2 my-12">
//         <label className="block text-sm font-medium text-secondary-700 mb-1">Enter Schedule</label>
//         <input
//           type="text"
//           name="cronSchedule"
//           value={data?.cronSchedule || ''}
//           onChange={handleChange}
//           placeholder="* 5 * * *"
//           className="w-full px-3 py-2 border border-secondary-300 rounded-md focus:outline-none focus:ring-1 focus:ring-primary-500"
//         />
//         {error && <span className="text-yellow-500 text-sm">{error}</span>}
//       </div>
//     </div>
//   );
// };

const Audit = () => {
  const postData = usePostAuditData();

  const handleComplete = (data) => {
    const selectedOptionEntry = Object.values(data).find((item) => item.selectedOption);
    const selectedOption = selectedOptionEntry?.selectedOption?.toLowerCase();
    const transformedData = transformResourceAuditPayload(data);
    const payload = { transformedData, selectedOption };
    postData.mutate(payload);
  };
  const steps = [
    {
      title: 'Cloud Account',
      component: (props) => <CreateCloud {...props} audit={true} />,
    },

    {
      title: 'Resource Audit',
      component: ResourceAudit,
    },
    // {
    //   title: 'Schedule',
    //   component: ScheduleStep,
    // },
  ];
  return (
    <div className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8 ">
      <Stepper steps={steps} handleComplete={handleComplete} postData={postData} />
    </div>
  );
};

export default Audit;
