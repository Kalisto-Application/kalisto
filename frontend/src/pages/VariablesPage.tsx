import React from 'react';
import { Environments } from '../components/Environments';
import { VariablesEditor } from '../components/VariablesEditor';
import VarsError from '../components/VarsError';
import { WorkspaceList } from '../components/Workspaces';
import SearchBox from '../ui/SearchBox';
import TabList from '../ui/TabList';

type VariablesPageProps = {};

export const VariablesPage: React.FC<VariablesPageProps> = () => {
  return (
    <div className="flex w-full flex-1">
      <div className="flex flex-[0_0_220px] flex-col justify-items-start">
        <WorkspaceList />
        <SearchBox />
        <Environments />
      </div>
      <div className="flex flex-1 flex-col">
        <TabList />
        <VariablesEditor />
        <div className="w-1/2"></div>
      </div>
      <VarsError />
    </div>
  );
};
