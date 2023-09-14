import React from 'react';

interface propsButton {
  text: string;
  className: string;
  disabled?: boolean;
  onClick: (e: React.SyntheticEvent) => void;
}

const Button: React.FC<propsButton> = ({
  text,
  onClick,
  className,
  disabled,
}) => {
  return (
    <button
      className={`rounded-md px-4 py-2 font-semibold leading-6 ${className} disabled:opacity-25`}
      onClick={onClick}
      disabled={disabled}
    >
      {text}
    </button>
  );
};

export default Button;
