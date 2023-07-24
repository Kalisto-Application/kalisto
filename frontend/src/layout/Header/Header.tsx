import React, { useContext } from "react";
import { Context } from "../../state";
import { WorkspaceList } from "../../components/Workspaces";

export const Header: React.FC = () => {
  const ctx = useContext(Context);

  return (<header className="h-[92px] select-none border-b border-solid border-layoutBorder bg-transparent">
    <WorkspaceList items={ctx.state.workspaceList} activeWorkspace={ctx.state.activeWorkspace} />
  </header>)
};
