import React, { useContext, useEffect, useState } from 'react';
import TreeView, {
  flattenTree,
  INodeRendererProps,
} from 'react-accessible-treeview';
import expandIcon from '../../assets/icons/expand.svg';
import addIcon from '../../assets/icons/plus.svg';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import InputAddItem from '../ui/InputAddItem';
import RequestList from '../ui/RequestList';

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

export type dataRequestList = {
  name: string;
};

export const MethodCollection: React.FC = () => {
  const ctx = useContext(Context);

  const services = ctx.state.activeWorkspace?.spec.services;
  const selectedItem = ctx.state.activeMethod?.fullName;

  const serviceNames = services?.map((it) => it.fullName) || [];

  const newServicesName = {
    name: '',
    children:
      services?.map((it, indx) => {
        return {
          name: it.displayName,

          children:
            it.methods.map((met) => {
              return {
                name: met.name,
                children: [{ name: '' }],
              };
            }) || [],
        };
      }) || [],
  };

  const files: dataRequestList[] = [];

  const data = flattenTree(newServicesName);

  const addItem = (value: string) => {
    files.push({ name: value });
  };

  return (
    <TreeView
      data={data}
      className="pl-4"
      onExpand={() => {}}
      nodeRenderer={({
        element,
        getNodeProps,
        level,
        isExpanded,
        handleExpand,
      }: INodeRendererProps) => {
        {
          return (
            <div
              {...getNodeProps({ onClick: handleExpand })}
              style={{
                paddingLeft: 30 * (level <= 2 ? level - 1 : 1.3),
                paddingRight: '30px',
                marginBottom: '10px',
                cursor: 'pointer',
              }}
            >
              {level <= 2 ? <ArrowIcon isOpen={isExpanded} /> : null}
              {element.name}
              {level === 3 ? (
                <div>
                  <InputAddItem
                    addItem={(value) => addItem(value)}
                    text="New Request"
                  />
                  <div
                    style={{
                      paddingLeft: 15 * (level - 1),
                    }}
                  >
                    <RequestList files={files} />
                  </div>
                </div>
              ) : null}
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
type propsArrowIcon = {
  isOpen: boolean;
};
const ArrowIcon: React.FC<propsArrowIcon> = ({ isOpen }) =>
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
