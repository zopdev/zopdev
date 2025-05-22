import React, { useState, useEffect, useCallback } from 'react';
import { CheckIcon, ChevronLeftIcon } from '@heroicons/react/24/outline';
import { ChevronRightIcon } from '@heroicons/react/20/solid';
import Button from '@/components/atom/Button/index.jsx';
import { useNavigate } from 'react-router-dom';
import ErrorComponent from '@/components/atom/ErrorComponent/index.jsx';
import { toast } from '@/components/molecules/Toast/index.jsx';

const Stepper = ({ steps, handleComplete, postData }) => {
  const [currentStep, setCurrentStep] = useState(0);
  const [stepData, setStepData] = useState({});
  const [stepStatus, setStepStatus] = useState(
    steps.map((_, index) => (index === 0 ? 'active' : 'incomplete')),
  );
  const [stepsComplete, setStepsComplete] = useState(steps.map(() => false));
  const navigate = useNavigate();

  useEffect(() => {
    if (postData?.isSuccess) {
      toast.success('Cloud Account Created Successfully.');
      navigate('/');
    }
  }, [postData]);

  useEffect(() => {
    const newStepStatus = steps.map((_, index) => {
      if (index < currentStep) return 'completed';
      if (index === currentStep) return 'active';
      return 'incomplete';
    });
    setStepStatus(newStepStatus);
  }, [currentStep, steps]);

  const setCurrentStepComplete = useCallback(
    (isComplete) => {
      setStepsComplete((prev) => {
        if (prev[currentStep] === isComplete) return prev;
        const updated = [...prev];
        updated[currentStep] = isComplete;
        return updated;
      });
    },
    [currentStep],
  );

  const updateStepData = (data) => {
    setStepData((prevData) => ({
      ...prevData,
      [currentStep]: data,
    }));
  };

  const isCurrentStepComplete = () => {
    return stepsComplete[currentStep];
  };
  const goToNextStep = () => {
    if (currentStep < steps.length - 1 && isCurrentStepComplete()) {
      setCurrentStep((prev) => prev + 1);
    } else if (currentStep === steps.length - 1 && isCurrentStepComplete()) {
      handleComplete(stepData);
    }
  };

  const goToPreviousStep = () => {
    if (currentStep > 0) {
      setCurrentStep((prev) => prev - 1);
    }
  };

  const renderDesktopStepIndicators = () => {
    return (
      <div className="flex items-center justify-center mb-8 w-full">
        {steps.map((step, index) => (
          <React.Fragment key={index}>
            <div
              onClick={() => {
                if (stepStatus[index] === 'completed' || index <= currentStep) {
                  setCurrentStep(index);
                }
              }}
              className={`
                flex items-center rounded-md px-4 py-2 text-center cursor-pointer
                ${
                  stepStatus[index] === 'completed'
                    ? 'border-2 border-primary-500'
                    : 'border border-primary-500'
                }
              `}
            >
              <div
                className={`
                  w-6 h-6 rounded-full flex items-center justify-center text-xs font-semibold mr-2
                  ${
                    stepStatus[index] === 'completed'
                      ? 'bg-primary-500 text-white'
                      : 'bg-secondary-300 text-secondary-800'
                  }
                `}
              >
                {stepStatus[index] === 'completed' ? <CheckIcon className="w-4 h-4" /> : index + 1}
              </div>
              <span className="text-sm font-medium text-secondary-800">{step.title}</span>
            </div>

            {index < steps.length - 1 && (
              <div
                className={`
                  flex-1 mx-2
                  ${stepStatus[index] === 'completed' ? 'h-[2px] bg-primary-500' : 'h-[2px] bg-secondary-300'}
                `}
              />
            )}
          </React.Fragment>
        ))}
      </div>
    );
  };

  const renderMobileDotIndicator = () => {
    return (
      <div className="flex flex-col items-center mb-6 w-full">
        <div className="flex items-center justify-center space-x-2 mb-2">
          {steps?.map((_, index) => (
            <div
              key={index}
              onClick={() => {
                if (stepStatus[index] === 'completed' || index <= currentStep) {
                  setCurrentStep(index);
                }
              }}
              className={`
                w-3 h-3 rounded-full cursor-pointer
                ${
                  index === currentStep
                    ? 'bg-primary-500'
                    : stepStatus[index] === 'completed'
                      ? 'bg-primary-300'
                      : 'bg-secondary-300'
                }
              `}
            />
          ))}
        </div>
        <div className="text-sm font-medium">
          Step {currentStep + 1}: {steps[currentStep].title}
        </div>
      </div>
    );
  };

  const renderCurrentStep = () => {
    const CurrentStepComponent = steps[currentStep].component;
    return (
      <CurrentStepComponent
        data={stepData[currentStep] || {}}
        updateData={updateStepData}
        allData={stepData}
        isLastStep={currentStep === steps.length - 1}
        setIsComplete={setCurrentStepComplete}
      />
    );
  };

  return (
    <div className="w-full  p-4 md:p-6 bg-white">
      <div className="hidden md:block">{renderDesktopStepIndicators()}</div>
      <div className="block md:hidden">{renderMobileDotIndicator()}</div>

      <div className="mb-8">
        {/* <h2 className="text-xl font-semibold mb-4">{steps[currentStep].title}</h2> */}
        {renderCurrentStep()}
      </div>
      <div className={' my-5'}>
        {postData?.isError && <ErrorComponent errorText={postData?.error?.message} />}
      </div>

      <div className="flex justify-between">
        {currentStep !== 0 ? (
          <Button
            variant="secondary"
            className="flex items-center px-3 py-2 md:px-4 md:py-2 rounded-md bg-secondary-200 text-secondary-700 hover:bg-secondary-300"
            onClick={goToPreviousStep}
            startEndornment={<ChevronLeftIcon className="w-3 h-3 mr-1" />}
          >
            {'Back'}
          </Button>
        ) : (
          <div />
        )}

        <Button
          loading={postData?.isPending}
          className={`flex items-center px-3 py-2 md:px-4 md:py-2 rounded-md ${
            !isCurrentStepComplete()
              ? 'bg-primary-300 text-white cursor-not-allowed'
              : 'bg-primary-500 text-white hover:bg-primary-600'
          }`}
          disabled={!isCurrentStepComplete()}
          onClick={goToNextStep}
          endEndornment={
            currentStep !== steps?.length - 1 && <ChevronRightIcon className="w-4 h-4 ml-1" />
          }
        >
          {currentStep === steps?.length - 1 ? 'Finish' : 'Next'}
        </Button>
      </div>
    </div>
  );
};

export default Stepper;
