import React, { useState } from "react";

interface props {
  items: itemProps[];
}

export const Menu: React.FC<props> = ({ items }) => {
  return (
    <div className="flex justify-center items-center px-0 py-1 flex-col content-normal rounded-md border-[1px] border-borderFill absolute z-10 bg-primaryFill cursor-pointer">
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
  const [isHovered, setIsHovered] = useState(false);

  const onMouseEnter = () => {
    setIsHovered(true);
  };

  const onMouseLeave = () => {
    setIsHovered(false);
  };

  let className =
    "flex flex-row w-[259.5px] h-11 items-center gap-[42px] px-4 py-2.5 leading-6";
  if (isHovered) {
    className += " bg-textBlockFill";
  }

  return (
    <div
      className={className}
      onClick={onClick}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      {icon && <img src={icon} />}
      <div>
        <div>{text}</div>
      </div>
    </div>
  );
};
