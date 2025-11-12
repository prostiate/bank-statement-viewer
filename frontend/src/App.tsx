import "./App.css";
import CsvUploader from "./components/CsvUploader";

function App() {
  return (
    <>
      <h1>Bank Statement Viewer</h1>
      <div className="card">
        <CsvUploader />
      </div>
    </>
  );
}

export default App;
