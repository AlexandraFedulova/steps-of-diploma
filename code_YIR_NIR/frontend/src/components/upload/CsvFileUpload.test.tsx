import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { CsvFileUpload } from "./CsvFileUpload";

describe("CsvFileUpload", () => {
  it("renders dropzone and file picker", () => {
    render(<CsvFileUpload file={null} onFileChange={vi.fn()} />);

    expect(screen.getByLabelText("Зона перетаскивания CSV")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Выбрать файл" })).toBeInTheDocument();
    expect(screen.getByText("Перетащите CSV сюда")).toBeInTheDocument();
  });

  it("accepts a CSV file via file input", async () => {
    const user = userEvent.setup();
    const onFileChange = vi.fn();

    render(<CsvFileUpload file={null} onFileChange={onFileChange} />);

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    const file = new File(["a,b\n1,0"], "matrix.csv", { type: "text/csv" });

    await user.upload(input, file);

    expect(onFileChange).toHaveBeenCalledWith(file);
  });

  it("rejects non-CSV files", () => {
    const onFileChange = vi.fn();

    render(<CsvFileUpload file={null} onFileChange={onFileChange} />);

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    const file = new File(["data"], "notes.txt", { type: "text/plain" });

    fireEvent.change(input, { target: { files: [file] } });

    expect(onFileChange).not.toHaveBeenCalled();
    expect(screen.getByRole("alert")).toHaveTextContent(/\.csv/i);
  });

  it("shows selected file name and clears on delete", async () => {
    const user = userEvent.setup();
    const onFileChange = vi.fn();
    const file = new File(["x"], "group.csv", { type: "text/csv" });

    const { rerender } = render(<CsvFileUpload file={null} onFileChange={onFileChange} />);

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    await user.upload(input, file);

    rerender(<CsvFileUpload file={file} onFileChange={onFileChange} />);

    expect(screen.getByText("group.csv")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Удалить" })).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Удалить" }));

    expect(onFileChange).toHaveBeenLastCalledWith(null);
  });

  it("does not accept files when disabled", async () => {
    const user = userEvent.setup();
    const onFileChange = vi.fn();

    render(<CsvFileUpload file={null} onFileChange={onFileChange} disabled />);

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    await user.upload(input, new File(["a"], "x.csv", { type: "text/csv" }));

    expect(onFileChange).not.toHaveBeenCalled();
    expect(screen.getByLabelText("Зона перетаскивания CSV")).toHaveAttribute("aria-disabled", "true");
  });
});
