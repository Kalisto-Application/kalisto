import React from "react";

import Header from "../Header";
import Sidebar from "../Sidebar";

type MainLayoutProps = {
  children?: React.ReactNode;
};

export const MainLayout: React.FC<MainLayoutProps> = ({ children }) => (
  <div className="flex h-screen flex-col">
    <Header>Header content</Header>
    <div className="flex flex-1">
      <Sidebar>Sidebar content</Sidebar>
      <main className="flex-1">{children}</main>
    </div>
  </div>
);
