import React from "react";

type SidebarProps = {

};

export const Sidebar: React.FC<SidebarProps> = ({}) => (
  <aside className="w-[80px] select-none border-r border-solid border-layoutBorder">
    <button>
      APIs
    </button>
    <button>
      Variables
    </button>
  </aside>
);
