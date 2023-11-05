import React, { useContext, useEffect, useState } from 'react';
import TreeView, {
  flattenTree,
  INodeRendererProps,
} from 'react-accessible-treeview';
import expandIcon from '../../assets/icons/expand.svg';
import addIcon from '../../assets/icons/plus.svg';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
const findMethod = (
  s: models.Service[] = [],
  serviceName: string,
  name: string
): models.Method | undefined => {
  for (const service of s) {
    if (service.fullName == serviceName) {
      for (const method of service.methods) {
        if (method.fullName == name) {
          return method;
        }
      }
    }
  }
};

type Data = {
  display: string;
  isMethod: boolean;
  parent?: string;
};

export const MethodCollection: React.FC = () => {
  const ctx = useContext(Context);

  const services = ctx.state.activeWorkspace?.spec.services;
  const selectedItem = ctx.state.activeMethod?.fullName;

  const serviceNames = services?.map((it) => it.fullName) || [];

  const newFolder = {
    name: '',
    children:
      services?.map((it) => {
        return {
          name: it.displayName,
          children:
            it.methods.map((met) => {
              return {
                name: met.name,
              };
            }) || [],
        };
      }) || [],
  };

  const data = flattenTree(newFolder);
  return (
    <TreeView
      data={data}
      className="pl-4"
      nodeRenderer={({
        element,
        getNodeProps,
        level,
        isBranch,
        isExpanded,
      }: INodeRendererProps) => {
        {
          return (
            <div
              {...getNodeProps()}
              style={{
                paddingLeft: 30 * (level - 1),
                marginBottom: '10px',
                cursor: 'pointer',
              }}
            >
              {isBranch ? (
                <FolderIcon isOpen={isExpanded} />
              ) : (
                <img src={addIcon} className="mr-2" />
              )}
              {element.name}
            </div>
          );
        }
      }}
    />
  );
};

/*
1) заменить либу и понять что работает ✅
2) научиться отображать фиксированный список в FileList
3) науиться создавать файлы
4) отоброжать активный файл на клик
5) подключить меню
*/
type propsFolder = {
  isOpen: boolean;
};
const FolderIcon: React.FC<propsFolder> = ({ isOpen }) =>
  isOpen ? (
    <img
      src={expandIcon}
      className="rotate-270 mr-2 transition duration-300 ease-in-out"
    />
  ) : (
    <img
      src={expandIcon}
      className="mr-2 -rotate-90 transition duration-300 ease-in-out"
    />
  );
