import React from "react";
import SearchBox from "../ui/SearchBox";
import { MethodCollection } from "../components/MethodCollectionView";
import { WorkspaceList } from "../components/Workspaces";
import { RequestEditor } from "../components/RequestEditor";
import { ResponseText } from "../components/ResponseText";
import TabList from "../ui/TabList";
import { ApiRequestSender } from "../components/ApiRequestSender";
import ApiError from "../components/ApiError";

export const ApiPage: React.FC = () => {
  return (
    <div className="flex flex-1 w-full">
      <div className="flex flex-[0_0_220px] justify-items-start flex-col">
        <WorkspaceList />
        <SearchBox />
        <MethodCollection />
      </div>
      <div className="flex flex-1 flex-col">
        <TabList />
        <ApiRequestSender />
        <div className="flex flex-1">
          <RequestEditor />
          <ResponseText />
        </div>
        <ApiError />
      </div>
    </div>
  );
};
