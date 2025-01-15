# `charsm` Documentation: Introducing `huh` Forms

The `huh` module in `charsm` is set of interactive components to create CLI-based forms using native shared libraries (`.dll` and `.so` files). its in an MVP stage.

---

## **Getting Started**

### **System Requirements**
1. **Node.js**: Ensure you have a Node.js version that supports native modules (`node-gyp`).
2. **Native Shared Library**: The module uses shared libraries (`.dll` for Windows, `.so` for Linux) to port Huh. Currently, macOS is **not supported(see contributions to build one)**.
3. **Build Instructions**:
   - To build the shared library:
     ```bash
     go build -o huh.dll -buildmode=c-shared .
     ```
     Replace `.dll` with `.so` for Linux systems.

Node has an option to setup native modules on install

![node](https://i.imgur.com/d7QDP8M.png)

---

## Good first pr 

in the fields.go file, combine `func (i *NewInput) Run() string ` and `func (i *NewInput) ForGroup() huh.Field` they do the same thing, Run should call forgroup and just run, group was an afterthought I got lazy!

make sure to sastify the interface 

```go
// remove absolete, value earguly returned on run
type Fieldsinterface interface {
	GetValue() string
	Run() string
	ForGroup() huh.Field
}
```

## **Features**

### **1. Set a Theme**
Customize the appearance of your form with pre-defined themes. 

```javascript
huh.SetTheme("dracula");
```

---

### **2. Create a Confirmation Dialog**
Create a simple confirmation dialog with customizable labels for Yes and No buttons.

```javascript
const m = huh.Confirm("Do you want to proceed?", "Yes", "No");
// To display the user's selection:
console.log(m.value);
```

---

### **3. Create an Input Field**
Define input fields with validation and placeholders.

#### Example 1: Single Input
```javascript
const i = new huh.NewInput(
  {
    Title: "Name",
    Description: "Enter your name",
    Placeholder: "e.g., John Doe",
    validators: "no_numbers,required",
  },
  0 // Mode: Single Input
);
i.load();
console.log(i.run());
// console.log(i.value);
```

#### Example 2: Multiline Input
```javascript
const s = new huh.NewInput(
  {
    Title: "Search",
    Description: "Enter query",
    Placeholder: "Type something...",
    validators: "required",
  },
  1 // Mode: Multiline
);
s.load();
console.log(s.run())
// console.log(s.value);
```

---

### **4. Create a Selection Component**
Allow users to select an option from a list.

```javascript
const c = huh.Select("Choose your favorite person", ["Opt1", "Opt2", "Opt2"]);
console.log(c.value);
```

---

### **5. Add Notes: DO NOT USE(BUGGY)**
Attach a note to the form (currently bugged).

```javascript
const note = huh.Note("Reminder", "This is a note", "NEXT LABEL", true);
// note.run();
```

---

### **6. Group Components**
Group multiple components together to form a cohesive form.

```javascript
const g1 = huh.CreateGroup(`${note.id},${i.id},${s.id},${m.id},${c.id}`);
huh.CreateForm(g1);
// GET VALUES E.G
console.log(i.value)
```

---

### **7. Add a Spinner**
Display a spinner to indicate loading or progress.

```javascript
huh.Spinner(2, "Preparing...");
```

---

## **Limitations and Notes**
1. **Known Bugs**:
   - The `huh.Note` and `huh.multiselect` components have issues.
2. **Native Module Requirement**:
   - Requires Node.js with `node-gyp` support to run shared libraries. Install necessary build tools:
     ```bash
     npm install -g node-gyp
     ```
   - Ensure GCC and `make` are installed (Linux) or Visual Studio Build Tools (Windows).
3. **Rapid Changes**:
   - This is an MVP (Minimum Viable Product), and the API may change rapidly in future versions.
4. **MacOS Compatibility**:
   - macOS users need to clone the repository, build the shared library, and test manually. Submit a PR with your changes:
     ```bash
     go build -o huh.dylib -buildmode=c-shared .
     ```

---

## Node.js with Native Build Tools(YOU NEED THIS TO RUN SHARED/NATIVE C/C++ LIBS)

On Windows, this is often referred to as installing Node.js with the **`node-gyp` prerequisites**, which include:

1. **Python** (Version 3.x recommended)
2. **Visual Studio Build Tools** or **MSVC (Microsoft Visual C++) Compiler** for Windows
3. **`node-gyp`**, which is a tool for building native modules in Node.js.

On Linux and macOS, it typically involves ensuring you have GCC, Make, and other development tools installed.

---

### **Installation Steps**

#### **Windows**
1. **Install Node.js**:
   Download the latest LTS version of Node.js from [nodejs.org](https://nodejs.org).

2. **Select native support in the Node Setup it'll handle everything**

4. **Install `node-gyp` globally**:
   Run the following command in PowerShell or Command Prompt:
   ```bash
   npm install -g node-gyp
   ```

---

#### **Linux**
1. **Install Build Essentials**:
   ```bash
   sudo apt update
   sudo apt install build-essential gcc g++ make python3
   ```

2. **Install `node-gyp` globally**:
   ```bash
   npm install -g node-gyp
   ```

---

#### **macOS**
1. **Install Xcode Command Line Tools**:
   ```bash
   xcode-select --install
   ```

2. **Install Python** (if not already installed):
   - Use Homebrew: `brew install python`.

3. **Install `node-gyp` globally**:
   ```bash
   npm install -g node-gyp
   ```

---

### **Validation**
After setup, validate the environment by creating a simple native module or running:
```bash
node-gyp configure build
```

This ensures that your Node.js environment is ready for compiling native modules.


## **Example Workflow**
Here’s a complete example that combines all components:

```javascript
import { huh } from "charsm";

huh.SetTheme("dracula");

const input = new huh.NewInput(
  {
    Title: "Name",
    Description: "Enter your name",
    Placeholder: "e.g., John Doe",
    validators: "required",
  },
  0
);
input.load();
input.run();

const confirmation = huh.Confirm("Do you want to proceed?", "Yes", "No");

const selection = huh.Select("Choose your option", ["Option 1", "Option 2"]);

const spinner = huh.Spinner(3, "Processing...");

console.log("Input Value:", input.value);
console.log("Confirmation Value:", confirmation.value);
console.log("Selection Value:", selection.value);
```

---

## **Feedback and Contribution**
- **Bugs**: Please report bugs via the issue tracker.
- **Pull Requests**: Contributions are welcome. For macOS compatibility, submit PRs with build/test updates.

