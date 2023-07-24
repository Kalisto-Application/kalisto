import React from "react";

import { Link } from "react-router-dom";

type SidebarProps = {

};

export const Sidebar: React.FC<SidebarProps> = ({}) => (
  <aside className="w-[80px] select-none border-r border-solid border-layoutBorder">
    <Link to="/api">APIs</Link>
    <Link to="/variables">Variables</Link>
  </aside>
);
