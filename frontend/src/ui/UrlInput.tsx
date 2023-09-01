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
    <div className="flex h-[50px]  border-2 border-borderFill">
      <input
        type="text"
        value={value}
        onChange={handleInputChange}
        className=" w-full bg-textBlockFill text-secondaryText"
      ></input>
      <button
        onClick={onClick}
        className=" flex-[0_0_106px] rounded-l-lg bg-primaryGeneral"
      >
        Send
      </button>
    </div>
  );
};
