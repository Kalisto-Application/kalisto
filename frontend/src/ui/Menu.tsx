import React, { useState } from 'react';

interface props {
  items: itemProps[];
}

export const Menu: React.FC<props> = ({ items }) => {
  return (
    <div className=" relative z-10  cursor-pointer rounded-md border-[1px] border-borderFill bg-primaryFill">
      {items.map((it, i) => {
        return <MenuItem key={i} {...it} />;
      })}
    </div>
  );
};

export default Menu;

interface itemProps {
  text: string;
  icon?: string;
  onClick?: (e: React.MouseEvent) => void;
}

export const MenuItem: React.FC<itemProps> = ({ text, icon, onClick }) => {
  return (
    <div
      className="flex gap-x-5 px-4 py-2.5 hover:bg-textBlockFill"
      onClick={onClick}
    >
      {icon && <img src={icon} />}
      <span>{text}</span>
    </div>
  );
};
