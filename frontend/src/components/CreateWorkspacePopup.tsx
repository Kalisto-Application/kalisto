import React, { useContext, useState } from 'react';
import { FindProtoFiles, CreateWorkspace } from './../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';
import { Context } from '../state';
import Button from '../ui/Button';
import Popup from '../ui/Popup';
import upload from '../../assets/icons/upload.svg';
import folder from '../../assets/icons/folder.svg';
import close from '../../assets/icons/close.svg';
import file from '../../assets/icons/file.svg';

interface propsCreateWorkspacePopup {
  onClose: () => void;
  open: boolean;
}
interface propsCreateWorkspace {
  onClose: () => void;
}

interface propsUpload {
  data: models.ProtoDir;

  setData: (dir: models.ProtoDir) => void;
  clearState: () => void;
}

const CreateWorkspacePopup: React.FC<propsCreateWorkspacePopup> = ({
  onClose,
  open,
}) => {
  return (
<Popup onClose={onClose} isOpen={open} title="Create Workspace">
      <CreateWorkspaceComponets onClose={onClose} />
    </Popup>
  );
};

export default CreateWorkspacePopup;

const CreateWorkspaceComponets: React.FC<propsCreateWorkspace> = ({
  onClose,
}) => {
  const ctx = useContext(Context);
  const [valueInp, setValueInp] = useState('');

  const [data, setData] = useState(
    new models.ProtoDir({
      folder: '',
      files: [],
    })
  );

  const updateTextValueInp = (e: React.FormEvent<HTMLInputElement>) => {
    let newValue = e.currentTarget.value;
    setValueInp(newValue);
  };

  const createNewWorkspace = () => {
    CreateWorkspace(valueInp, data.folder).then((res) => {
      ctx.dispatch({ type: 'newWorkspace', workspace: res });
      onClose();
    });
  };

  return (
    <div className="grid">
      <div className="mb-10 grid gap-y-2">
        <Upload
          data={data}
          setData={setData}
          clearState={() =>
            setData({
              folder: '',
              files: [],
            })
          }
        />
        <label className="leading-6">Name</label>
        <input
          value={valueInp}
          onChange={(e) => updateTextValueInp(e)}
          className="rounded-md border-[1px] border-borderFill bg-primaryFill px-4 py-1.5 placeholder-secondaryText"
          type="text"
          placeholder="Name of Workspace"
        />
      </div>
      <div className="justify-self-end">
        <Button
          text="Cancel"
          onClick={onClose}
          className="mr-[10px] border-[1px] border-[#343434]"
        />
        <Button
          disabled={valueInp === '' || data.files.length <= 0}
          text="Create"
          onClick={() => createNewWorkspace()}
          className="bg-primaryGeneral text-lg font-medium"
        />
      </div>
    </div>
  );
};

const Upload: React.FC<propsUpload> = ({ data, clearState, setData }) => {
  const [disabledUpload, setDisabledUpload] = useState(false);

  const findFiles = () => {
    setDisabledUpload(true);
    FindProtoFiles()
      .then((res) => {
        if (res.files.length !== 0 && res.folder !== '') {
          setData({
            folder: res.folder,
            files: res.files,
          });
        }
        setDisabledUpload(false);
      })
      .catch((er) => setDisabledUpload(false));
  };

  return (
    <div>
      {data.files.length === 0 ? (
        <button
          onClick={findFiles}
          className="mb-6 flex min-h-[179px] w-full flex-col items-center justify-center border-2 border-dashed border-borderFill"
          disabled={disabledUpload}
        >
          <img src={upload} className="mb-6" alt="upload" />
          <span className="font-sm  font-['Roboto_Mono'] font-bold text-secondaryText">
            Choose folder to upload
          </span>
          <span className="text-xs  text-blind">proto</span>
        </button>
      ) : (
        <div className="max mb-6 max-h-52 overflow-auto rounded-md bg-textBlockFill px-5 py-3.5 ">
          <div className="flex justify-between [&:not(:last-child)]:mb-1">
            <div className="flex items-center gap-x-2">
              <img src={folder} alt="file" />
              <span className="text-sm text-blueTextPath">{data.folder}</span>
            </div>
            <button
              onClick={() => {
                clearState();
              }}
            >
              <img src={close} alt="clearState" />
            </button>
          </div>
          <ul>
            {data.files.map((f, indx) => {
              return (
                <li
                  key={indx}
                  className="flex gap-x-2  pl-5 text-sm [&:not(:last-child)]:mb-1"
                >
                  <img src={file} alt="file" />
                  <span>{f}</span>
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </div>
  );
};
