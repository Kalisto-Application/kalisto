import { useState } from 'react';
import { useBoolean } from 'usehooks-ts';
import plusIcon from '../../assets/icons/plus.svg';

type inputProps = {
  text: string;
  addItem: (value: string) => void;
  placeholder?: string;
};
const InputAddItem: React.FC<inputProps> = ({ addItem, text, placeholder }) => {
  const {
    value: showInput,
    setValue: setShowInput,
    setTrue,
    setFalse,
  } = useBoolean(false);

  const [inputValue, setInputValue] = useState('');
  const [valueValidate, setValueValidate] = useState(false);

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.code == 'Enter' && inputValue !== '') {
      setValueValidate(false);
      addItem(inputValue);
      setFalse();
      setInputValue('');
    }
  };

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let value = e.target.value;
    setValueValidate(false);
    setInputValue(value);
  };

  return (
    <div className="mb-3 flex flex-col">
      {showInput ? (
        <>
          <input
            autoFocus
            onKeyDown={onKeyDown}
            type="text"
            placeholder={placeholder}
            onChange={(e) => onChange(e)}
            className="border-1 mb-1 w-full border-[1px] border-borderFill  bg-textBlockFill px-3 placeholder:text-[14px] placeholder:text-secondaryText"
          />

          <div className="pl-5 text-xs">Push Enter to rename </div>
        </>
      ) : (
        <>
          <button
            className="flex items-center gap-x-2  rounded-md border-borderFill bg-primaryFill  transition duration-500 ease-in-out hover:bg-textBlockFill"
            onClick={setTrue}
          >
            <img src={plusIcon} alt="" />
            <span className="text-lg">{text}</span>
          </button>
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
export default InputAddItem;
