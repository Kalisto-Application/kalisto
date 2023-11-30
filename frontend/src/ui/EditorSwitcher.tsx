import React, { useContext, useRef, useState } from 'react';
import copyIcon from '../../assets/icons/copy.svg';
import PopoverUI from './Popover';

interface props {
  items: item[];
  active: number;
  onClickCopy: () => void;
}

type item = {
  title: string;
  onClick: (i: number) => void;
};

export const EditorSwitcher: React.FC<props> = ({
  items,
  active,
  onClickCopy,
}) => {
  const makeClassName = (i: number): string => {
    return active === i ? '' : 'text-secondaryText';
  };

  return (
    <div className="relative  z-10 flex  border-l-2 border-borderFill py-2">
      <div className="flex  flex-1  font-['Roboto_Mono']">
        {items.map((it, i) => (
          <div
            className={`${makeClassName(i)} cursor-pointer px-4`}
            key={i}
            onClick={() => it.onClick(i)}
          >
            {it.title}
          </div>
        ))}
      </div>
      <div className="mr-8 flex">
        <PopoverUI text="Copied">
          <button className="self-center" onClick={() => onClickCopy()}>
            <img src={copyIcon} />
          </button>
        </PopoverUI>
      </div>
    </div>
  );
};
