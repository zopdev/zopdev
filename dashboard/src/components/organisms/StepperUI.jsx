import Stepper from '@/components/organisms/Stepper.jsx';
import { useCreateResourceMutation } from '@/Queries/CloudAccount/index.js';
import { useNavigate } from 'react-router-dom';

export default function StepperUI({ steps }) {
  const postData = useCreateResourceMutation();
  const navigate = useNavigate();
  const handleComplete = (data) => {
    postData.mutate({
      data,
    });
    if (postData?.isSuccess) {
      navigate('/');
    }
  };

  return (
    <div className=" mx-auto py-10">
      <Stepper steps={steps} onComplete={handleComplete} />
    </div>
  );
}
