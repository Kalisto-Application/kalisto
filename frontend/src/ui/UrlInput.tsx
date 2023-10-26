import React, { SyntheticEvent } from 'react';

interface UrlInputProps {
  value: string;
  setValue: (v: string) => void;
  onClick: (e: SyntheticEvent) => void;
}

export const UrlInput: React.FC<UrlInputProps> = ({
  onClick,
  value,
  setValue,
}) => {
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
  };

  return (
    <div className="flex border-2 border-borderFill">
      <input
        type="text"
        value={value}
        onChange={handleInputChange}
        className="w-full bg-textBlockFill px-5 font-['Roboto_Mono'] text-secondaryText"
      ></input>
      <button
        onClick={onClick}
        className=" rounded-l-lg bg-primaryGeneral px-8 py-3"
      >
        Send
      </button>
    </div>
  );
};
