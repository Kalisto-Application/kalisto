import React, { useContext, useRef, useState } from 'react';
import copyIcon from '../../assets/icons/copy.svg';

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
    return active === i
      ? 'pt-[8px] px-[16px]'
      : 'pt-[8px] px-[16px] text-secondaryText ';
  };

  return (
    <div className="ustify-between flex">
      <div className="flex h-[40px] flex-1  ">
        {items.map((it, i) => (
          <div
            className={`${makeClassName(i)} cursor-pointer`}
            key={i}
            onClick={() => it.onClick(i)}
          >
            {it.title}
          </div>
        ))}
      </div>
      <button onClick={onClickCopy}>
        <img src={copyIcon} />
      </button>
    </div>
  );
};
