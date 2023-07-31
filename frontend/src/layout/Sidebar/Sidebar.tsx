import React, { useEffect } from "react";

import { Link, useMatch } from "react-router-dom";

export const Sidebar: React.FC = ({}) => (
  <aside className="w-[90px] flex flex-col justify-star border-[1px] border-r-borderFill border-y-0 border-l-0">
    <SideIcon link="/api" />
    <SideIcon link="/variables" />
  </aside>
);


type SideIconProps = {
  link: string;
}
const SideIcon: React.FC<SideIconProps> = ({ link }) => {
  let className = "h-[40px] my-[10px] py-[8px] px-[26px] static"
  const match = useMatch(link);
  if (match?.pathname === link) {
    className += " border-l-primaryGeneral border-l-4"
  } else {
    className += " border-l-primaryFill border-l-4"
  }

  return <div className={className}>
    <Link to={link}>
      <svg width="26" height="27" viewBox="0 0 26 27" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M7.34574 0.5H4.68221C1.61214 0.5 0 2.11215 0 5.16822V7.83178C0 10.8878 1.61214 12.5 4.66819 12.5H7.33172C10.3878 12.5 11.9999 10.8878 11.9999 7.83178V5.16822C12.0139 2.11215 10.4018 0.5 7.34574 0.5Z" fill="white"/>
        <path d="M21.189 0.5H18.5254C15.4693 0.5 13.8572 2.11215 13.8572 5.16822V7.83178C13.8572 10.8878 15.4693 12.5 18.5254 12.5H21.189C24.245 12.5 25.8572 10.8878 25.8572 7.83178V5.16822C25.8572 2.11215 24.245 0.5 21.189 0.5Z" fill="#3D3DAB"/>
        <path d="M21.189 14.3574H18.5254C15.4693 14.3574 13.8572 15.9696 13.8572 19.0256V21.6892C13.8572 24.7453 15.4693 26.3574 18.5254 26.3574H21.189C24.245 26.3574 25.8572 24.7453 25.8572 21.6892V19.0256C25.8572 15.9696 24.245 14.3574 21.189 14.3574Z" fill="white"/>
        <path d="M7.34574 14.3574H4.68221C1.61214 14.3574 0 15.9677 0 19.0202V21.6806C0 24.7472 1.61214 26.3574 4.66819 26.3574H7.33172C10.3878 26.3574 11.9999 24.7472 11.9999 21.6946V19.0342C12.0139 15.9677 10.4018 14.3574 7.34574 14.3574Z" fill="white"/>
      </svg>
    </Link>
  </div>
}