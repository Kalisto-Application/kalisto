import React from "react";

import { Link, useMatch } from "react-router-dom";

import apiIcon from "../../../assets/icons/sideApi.svg";
import varsIcon from "../../../assets/icons/sideVars.svg";
import scriptIcon from "../../../assets/icons/sideScript.svg";
import apiIconActive from "../../../assets/icons/sideApiActive.svg";
import varsIconActive from "../../../assets/icons/sideVarsActive.svg";
import scriptIconActive from "../../../assets/icons/sideScriptActive.svg";

export const Sidebar: React.FC = ({}) => (
  <aside className="justify-star flex flex-[0_0_90px] flex-col border-[1px] border-y-0 border-l-0 border-r-borderFill">
    <SideIcon link="/api" icon={apiIcon} activeIcon={apiIconActive} />
    <SideIcon
      link="/scripting"
      icon={scriptIcon}
      activeIcon={scriptIconActive}
    />
    <SideIcon link="/variables" icon={varsIcon} activeIcon={varsIconActive} />
  </aside>
);

type SideIconProps = {
  link: string;
  icon: string;
  activeIcon: string;
};

const SideIcon: React.FC<SideIconProps> = ({ link, icon, activeIcon }) => {
  // let className = 'h-[40px] w-full my-[10px] py-[8px] px-[26px] static';
  let iconSrc = icon;
  const match = useMatch(link);
  if (match?.pathname === link) {
    iconSrc = activeIcon;
    // className += ' border-l-primaryGeneral border-l-4';
  } else {
    // className += ' border-l-primaryFill border-l-4';
  }

  return (
    <Link to={link}>
      <img src={iconSrc} />
    </Link>
  );
};
