import React, { useState } from 'react';
import Button from '../ui/Button';
import Popup from '../ui/Popup';
import upload from '../../assets/icons/upload.svg';
import folder from '../../assets/icons/folder.svg';
import close from '../../assets/icons/close.svg';
import file from '../../assets/icons/file.svg';

// State для Popup
//   const [isOpenCreateWorkspace, setIsOpenCreateWorkspace] = useState(false);

// кнопка для вызова Popup

// <button onClick={() => setIsOpenDeletePopup(true)}>open delete</button>;

// компонента для Popup
{
  /* <CreateWorkspacePopup
        open={isOpenCreateWorkspace}
        onClose={() => setIsOpenCreateWorkspace(false)}
      /> */
}

interface propsCreateWorkspacePopup {
  onClose: () => void;
  open: boolean;
}
interface propsCreateWorkspace {
  onClose: () => void;
}

interface propsUpload {
  data: {
    folder: string;
    files: string[];
  };
  clearState: () => void;
}

const CreateWorkspacePopup: React.FC<propsCreateWorkspacePopup> = ({
  onClose,
  open,
}) => {
  return (
    <Popup onClose={onClose} isOpen={open} title="Name of New Request">
      <CreateWorkspace onClose={onClose} />
    </Popup>
  );
};

export default CreateWorkspacePopup;

const CreateWorkspace: React.FC<propsCreateWorkspace> = ({ onClose }) => {
  const [data, setData] = useState({
    folder: '/my/folder',
    files: ['ast.proto', 'met.proto', 'super.proto'],
  });
  return (
    <div className="grid">
      <div className="mb-10 grid gap-y-2">
        <Upload
          data={data}
          clearState={() =>
            setData({
              folder: '',
              files: [],
            })
          }
        />
        <label className="leading-6">Name</label>
        <input
          className="rounded-md border-[1px] border-borderFill bg-primaryFill px-4 py-1.5 placeholder-secondaryText"
          type="text"
          placeholder="Name of New Request"
        />
      </div>
      <div className="justify-self-end">
        <Button
          text="Cancel"
          onClick={onClose}
          className="mr-[10px] border-[1px] border-[#343434]"
        />
        <Button
          text="Create"
          onClick={() => {}}
          className="bg-primaryGeneral text-lg font-medium"
        />
      </div>
    </div>
  );
};

const Upload: React.FC<propsUpload> = ({ data, clearState }) => {
  return (
    <div>
      <button className="mb-6 flex min-h-[179px] w-full flex-col items-center justify-center border-2 border-dashed border-borderFill">
        <img src={upload} className="mb-6" alt="upload" />
        <span className="font-sm  font-['Roboto_Mono'] font-bold text-secondaryText">
          Choose folder to upload
        </span>
        <span className="text-xs  text-[#5E5E5E]">proto</span>
      </button>
      {data.files.length !== 0 ? (
        <div className="mb-6 rounded-md bg-textBlockFill px-5 py-3.5">
          <div className="flex justify-between [&:not(:last-child)]:mb-1">
            <div className="flex items-center gap-x-2">
              <img src={folder} alt="file" />
              <span className="text-sm text-blueTextPath">{data.folder}</span>
            </div>
            <button onClick={clearState}>
              <img src={close} alt="clearState" />
            </button>
          </div>
          <ul>
            {data.files.map((f) => {
              return (
                <li className="flex gap-x-2  pl-5 text-sm [&:not(:last-child)]:mb-1">
                  <img src={file} alt="file" />
                  <span>{f}</span>
                </li>
              );
            })}
          </ul>
        </div>
      ) : null}
    </div>
  );
};
