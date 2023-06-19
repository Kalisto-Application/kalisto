import React, {SyntheticEvent, useState} from "react";

interface UrlInputProps {
  onClick: (e: SyntheticEvent, data: string) => void;
};

export const UrlInput: React.FC<UrlInputProps> = ({ onClick }) => {
  const [data, setData] = useState<string>('');

  const handleOnClick = (e: SyntheticEvent) => {
    onClick(e, data)
  }
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setData(e.target.value)
  }

  return(
  <div className="h-49 flex flex-1 justify-between bg-codeSectionBg px-2 py-1 rounded-t-md">
    <input type="text" value={data} onChange={handleInputChange} className="bg-transparent flex-1 text-inputPrimary mx-4 my-2"></input>
    <button onClick={handleOnClick} className="px-[33px] h-[27px] rounded-lg bg-buttonPrimary text-white mx-4 my-2">Send</button>
  </div>
  );
};
 