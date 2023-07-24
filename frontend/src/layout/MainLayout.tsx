import React from "react";

import {
  Outlet,
} from "react-router-dom";

import Header from "./Header";
import Sidebar from "./Sidebar";

export const MainLayout: React.FC = () => {
  return (<div className="flex h-screen flex-col">
    <Header />
    <div className="flex flex-1">
      <Sidebar />
      <Outlet />
    </div>
  </div>);
}
