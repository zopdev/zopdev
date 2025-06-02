function Checkbox({ children, className, ...props }) {
  return (
    <input
      type="checkbox"
      className={`h-4 w-4 accent-primary-600  rounded !bg-white  focus:accent-primary-500  focus:ring-offset-white ${className}`}
      {...props}
    />
  );
}

export default Checkbox;
