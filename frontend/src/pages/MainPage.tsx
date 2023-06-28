import React, {useState} from "react";
import MainLayout from "../layout/MainLayout";
import { UrlInput } from "../components/UrlInput";
import { CodeEditor } from "../components/CodeEditor";
import { MethodCollection, MethodItem } from "../components/MethodCollectionView";
import {SendGrpc} from "../../wailsjs/go/api/Api"

export const MainPage: React.FC = () => {
  return (
    <MainLayout />
  );
};
