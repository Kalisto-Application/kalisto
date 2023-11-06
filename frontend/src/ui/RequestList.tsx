import { dataRequestList } from '../components/MethodCollectionView';

type RequestList = {
  files: dataRequestList[];
};
const RequestList: React.FC<RequestList> = ({ files }) => {
  return (
    <>
      {files.map((it, indx) => (
        <div className="grid" key={indx}>
          {it.name}
        </div>
      ))}
    </>
  );
};
export default RequestList;
