import React from "react";

import { Link, useMatch } from "react-router-dom";

import apiIcon from "../../../assets/icons/sideApi.svg";
import varsIcon from "../../../assets/icons/sideVars.svg";
import scriptIcon from "../../../assets/icons/sideScript.svg";

export const Sidebar: React.FC = ({}) => (
  <aside className="flex flex-[0_0_90px] flex-col justify-star border-[1px] border-r-borderFill border-y-0 border-l-0">
    <SideIcon link="/api" icon={apiIcon} />
    <SideIcon link="/variables" icon={varsIcon} />
    <SideIcon link="/scripting" icon={scriptIcon} />
  </aside>
);


type SideIconProps = {
  link: string;
  icon: string;
}

const SideIcon: React.FC<SideIconProps> = ({ link, icon }) => {
  let className = "h-[40px] w-full my-[10px] py-[8px] px-[26px] static"
  const match = useMatch(link);
  if (match?.pathname === link) {
    className += " border-l-primaryGeneral border-l-4"
  } else {
    className += " border-l-primaryFill border-l-4"
  }

  return <div className={className}>
    <Link to={link}>
      <img src={icon}/>
    </Link>
  </div>
}