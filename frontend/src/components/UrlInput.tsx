import React, {SyntheticEvent} from "react";

interface UrlInputProps {
  value: string;
  setValue: (v: string) => void;
  onClick: (e: SyntheticEvent) => void;
};

export const UrlInput: React.FC<UrlInputProps> = ({ onClick, value, setValue }) => {
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value)
  }

  return(
  <div className="h-49 flex flex-1 justify-between bg-codeSectionBg px-2 py-1 rounded-t-md">
    <input type="text" value={value} onChange={handleInputChange} className="bg-transparent flex-1 text-inputPrimary mx-4 my-2"></input>
    <button onClick={onClick} className="px-[33px] h-[27px] rounded-lg bg-buttonPrimary text-white mx-4 my-2">Send</button>
  </div>
  );
};
