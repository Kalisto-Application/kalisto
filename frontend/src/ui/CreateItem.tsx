import { useState } from 'react';
import { useBoolean } from 'usehooks-ts';
import plusIcon from '../../assets/icons/plus.svg';

type CreateItem = {
  text: string;
  addItem: (value: string, fullNameMet?: string) => void;
  placeholder?: string;
  fullNameMet?: string;
  btnClassName?: string;
};

const CreateItem: React.FC<CreateItem> = ({
  addItem,
  text,
  placeholder,
  fullNameMet,
  btnClassName,
}) => {
  const { value: showInput, setTrue, setFalse } = useBoolean(false);

  const [inputValue, setInputValue] = useState('');
  const [valueValidate, setValueValidate] = useState(false);

  const onKeyDown = (e: React.KeyboardEvent): void => {
    if (e.code === 'Enter' && inputValue !== '') {
      setValueValidate(false);
      addItem(inputValue, fullNameMet);
      setFalse();
      setInputValue('');
    }
  };

  const onChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const value = e.target.value;
    setValueValidate(false);
    setInputValue(value);
  };

  const btnClass = `flex w-full items-center gap-x-2 border-borderFill bg-primaryFill hover:bg-borderFill ${
    btnClassName !== '' ? btnClassName : ''
  }`;

  return (
    <div className="flex w-full flex-col">
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

          <div className="pl-5 text-xs">Push Enter to save </div>
        </>
      ) : (
        <>
          <button className={btnClass} onClick={setTrue}>
            <img src={plusIcon} alt="" />
            <span className="whitespace-nowrap text-lg">{text}</span>
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
export default CreateItem;
