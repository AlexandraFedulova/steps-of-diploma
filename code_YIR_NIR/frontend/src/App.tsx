import { Route, Routes } from "react-router-dom";
import { AppLayout } from "@/components/layout/AppLayout";
import { GraphPage } from "@/pages/GraphPage";
import { HomePage } from "@/pages/HomePage";
import { UploadPage } from "@/pages/UploadPage";

export default function App() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route index element={<HomePage />} />
        <Route path="upload" element={<UploadPage />} />
        <Route path="graph" element={<GraphPage />} />
      </Route>
    </Routes>
  );
}
