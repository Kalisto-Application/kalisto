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
  <div className="h-[52px] flex flex-1 border-2 border-borderFill">
    <input type="text" value={value} onChange={handleInputChange} className="w-full h-full text-secondaryText bg-textBlockFill"></input>
    <button onClick={onClick} className="w-[106px] h-full rounded-l-lg bg-primaryGeneral">Send</button>
  </div>
  );
};
