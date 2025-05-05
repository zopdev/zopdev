function Label({ children, className, ...props }) {
    return (
        <label className={`block text-sm  leading-6 text-secondary-600 ${className}`} {...props}>
            {children}
        </label>
    );
}

export default Label;
