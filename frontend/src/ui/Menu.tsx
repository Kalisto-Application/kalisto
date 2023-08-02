import React, {useState} from "react";

interface props {
  items: itemProps[]
}

export const Menu: React.FC<props> = ({items}) => {
    return (
        <div className="border-[1px] border-borderFill z-10 bg-primaryFill p-2">
          {items.map((it, i) => {
            return <MenuItem key={i} {...it} />
          })}
        </div>
    )
}

export default Menu;

interface itemProps {
  text: string;
  icon?: string;
  onClick?: () => void;
};

export const MenuItem: React.FC<itemProps> = ({text, icon, onClick}) => {
  const [isHovered, setIsHovered] = useState(false);

  const onMouseEnter = () => {
    setIsHovered(true);
  }

  const onMouseLeave = () => {
    setIsHovered(false);
  }

  let className = "m-2";
  if (isHovered) {className += " bg-textBlockFill"}

    return (
        <div className={className} onClick={onClick} onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}>
          {icon && <img src={icon} />}
          <div>
            <div>{text}</div>
          </div>
        </div>
    )
}
