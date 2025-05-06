import Stepper from '@/components/organisms/Stepper.jsx';

export default function StepperUI({ steps }) {
  return (
    <div className=" mx-auto py-10">
      <Stepper steps={steps} />
    </div>
  );
}
