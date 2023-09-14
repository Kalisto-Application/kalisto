import React from "react";
import SearchBox from "../ui/SearchBox";
import { WorkspaceList } from "../components/Workspaces";
import { ScriptEditor } from "../components/ScriptEditor";
import { ScriptResponse } from "../components/ScriptResponse";
import TabList from "../ui/TabList";
import ScriptError from "../components/ScriptError";
import { ScriptSender } from "../components/ScriptSender";

export const ScriptingPage: React.FC = () => {
  return (
    <div className="flex flex-1 w-full">
      <div className="flex flex-[0_0_220px] justify-items-start flex-col">
        <WorkspaceList />
        <SearchBox />
      </div>
      <div className="flex flex-1 flex-col">
        <TabList />
        <ScriptSender />
        <div className="flex flex-1">
          <ScriptEditor />
          <ScriptResponse />
        </div>
        <ScriptError />
      </div>
    </div>
  );
};
