import React from 'react';

interface propsButton {
  text: string;
  className: string;
  onClick: (e: React.SyntheticEvent) => void;
}

const Button: React.FC<propsButton> = ({ text, onClick, className }) => {
  return (
    <button
      className={`rounded-md px-4 py-2 font-semibold leading-6 ${className}`}
      onClick={onClick}
    >
      {text}
    </button>
  );
};

export default Button;
