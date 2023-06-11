import React from "react";

type SidebarProps = {
  children?: React.ReactNode;
};

export const Sidebar: React.FC<SidebarProps> = ({ children }) => (
  <aside className="w-[90px] select-none border-r border-solid border-layoutBorder">
    {children}
  </aside>
);
