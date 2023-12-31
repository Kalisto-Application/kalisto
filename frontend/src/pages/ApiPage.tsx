import React from 'react';
import ApiError from '../components/ApiError';
import { ApiRequestSender } from '../components/ApiRequestSender';
import { MethodCollection } from '../components/MethodCollectionView';
import { RequestEditor } from '../components/RequestEditor';
import { ResponseText } from '../components/ResponseText';
import { WorkspaceList } from '../components/Workspaces';
import SearchBox from '../ui/SearchBox';
import TabList from '../ui/TabList';

export const ApiPage: React.FC = () => {
  return (
    <div className="flex w-full flex-1">
      <div className="flex flex-[0_0_220px] flex-col justify-items-start">
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
