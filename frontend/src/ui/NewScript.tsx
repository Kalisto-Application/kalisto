import { useState } from 'react';
import plusIcon from '../../assets/icons/plus.svg';

type ScriptNewProps = {
  addFile: (value: string) => void;
};

const NewScript: React.FC<ScriptNewProps> = ({ addFile }) => {
  const [value, setValue] = useState('');
  const [isMode, setIsMode] = useState(true);
  const [valueValidate, setValueValidate] = useState(false);

  const updateValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setValueValidate(false);
    setValue(newValue);
  };

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.code == 'Enter' && value === '') {
      setValueValidate(true);
      return;
    }

    if (e.code == 'Enter') {
      addFile(value);
      setIsMode(true);
      setValue('');
    }
  };

  const onClick = () => {
    setIsMode(false);
  };
  return (
    <div className="mb-3 flex flex-col">
      {isMode ? (
        <button
          className="flex items-center gap-x-2  rounded-md border-borderFill bg-primaryFill px-3 py-1 transition duration-500 ease-in-out hover:bg-textBlockFill"
          onClick={onClick}
        >
          <img src={plusIcon} alt="" />
          <span className="text-lg">Add Script</span>
        </button>
      ) : (
        <>
          <div className="px-5">
            <input
              placeholder="Name script"
              className="border-1 mb-1 w-full border-[1px] border-borderFill  bg-textBlockFill px-3 placeholder:text-[14px] placeholder:text-secondaryText"
              type="text"
              autoFocus
              onChange={(e) => updateValue(e)}
              onKeyDown={onKeyDown}
            />
          </div>
          <div className="pl-5 text-xs">Push Enter to rename </div>
        </>
      )}
      {valueValidate ? (
        <span className="text-[13px] text-red">
          A script name must not be empty
        </span>
      ) : null}
    </div>
  );
};

export default NewScript;
