document.addEventListener("DOMContentLoaded", function () {
  const log = (..._data) => {
    //console.log(..._data);
  };

  const urls = document.getElementById("editor~urls").dataset;
  const newBlockModal = document.getElementById("new-block~modal");

  const getFormData = (formElem, deep = true) => {
    const blockName = formElem.dataset.name;
    const blockForm = new FormData(formElem);
    var blockOpts = {};
    blockForm.forEach((v, k) => {
      blockOpts[k] = v;
    });
    var ret = {
      name: blockName,
      opts: blockOpts,
    };
    if (deep) {
      ret.children = Array.from(
        formElem.parentNode.querySelectorAll(formElem.dataset.childrenSelector),
      ).map(getFormData);
    }
    return ret;
  };

  const formEvents = (root) => {
    // children toggle
    const doChildToggle = (target, checked) => {
      if (checked) {
        target.style.display = "block";
      } else {
        target.style.display = "none";
      }
    };
    const listenChildToggle = (checkbox) => {
      log();
      log(checkbox.dataset);
      const formChildren = document.getElementById(checkbox.dataset.formTarget);
      log(formChildren);
      const block = document.getElementById(checkbox.dataset.blockTarget);
      log(block);
      const blockChildren = block.querySelector(block.dataset.childrenSelector);
      log(block.dataset, blockChildren);
      if (blockChildren != null) {
        checkbox.addEventListener("change", (e) => {
          doChildToggle(formChildren, e.target.checked);
        });
        doChildToggle(formChildren, checkbox.checked);
      } else {
        formChildren.remove();
        checkbox.parentNode.remove();
      }
    };
    const childToggles = root.querySelectorAll(".block-form\\~toggle-children");
    log(childToggles);
    childToggles.forEach(listenChildToggle);

    // delete button
    const doDelete = (form, block) => {
      if (confirm("Are you sure you want to delete this block?")) {
        form.remove();
        block.remove();
      }
    };
    const listenDelete = (btn) => {
      const form = document.getElementById(btn.dataset.formTarget);
      const block = document.getElementById(btn.dataset.blockTarget);
      btn.addEventListener("click", (_) => doDelete(form, block));
    };
    const deleteBtns = root.querySelectorAll(".block-form\\~delete");
    log(deleteBtns);
    deleteBtns.forEach(listenDelete);

    // preview form changes
    const doFormChanged = (target, data) => {
      const req = new Request(urls.blockUpdate, {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
          "Content-Type": "application/json",
        },
      });
      const bc = target.querySelector(target.dataset.childrenSelector);
      var children = null;
      if (bc != null) {
        children = bc.innerHTML;
      }
      fetch(req)
        .then((resp) => resp.text())
        .then((html) => {
          target.innerHTML = html;

          const bbc = target.querySelector(target.dataset.childrenSelector);
          if (bbc != null) {
            bbc.innerHTML = children;
          }
        })
        .catch(console.error);
    };
    const listenFormChanged = (input) => {
      input.addEventListener("change", (e) => {
        const form = e.target.parentElement;
        const data = getFormData(form, false);
        const id = form.dataset.id;
        const target = document.getElementById(form.dataset.blockTarget);
        doFormChanged(target, data, id);
      });
    };
    const formFields = root.querySelectorAll(".form-field");
    formFields.forEach(listenFormChanged);

    // create new block
    const newBlockButtons = root.querySelectorAll(".block-form\\~new");
    log("newBlockButtons", newBlockButtons);
    newBlockButtons.forEach((btn) => {
      btn.addEventListener("click", (_) => {
        newBlockModal.dataset.blockTarget = btn.dataset.blockTarget;
        newBlockModal.dataset.formTarget = btn.dataset.formTarget;
        newBlockModal.style.display = "block";
      });
    });
  };

  const modalEvents = () => {
    const newBlockName = document.getElementById("new-block~name");
    const newBlockCreate = document.getElementById("new-block~create");
    const newBlockCancel = document.getElementById("new-block~cancel");
    const newBlockWkSpace = document.getElementById("new-block~wkspace");
    log(
      "newBlock",
      newBlockModal,
      newBlockName,
      newBlockCreate,
      newBlockCancel,
      newBlockWkSpace,
    );

    newBlockCancel.addEventListener("click", (_) => {
      newBlockModal.style.display = "none";
    });
    newBlockCreate.addEventListener("click", (_) => {
      const blockParent = document.getElementById(
        newBlockModal.dataset.blockTarget,
      );
      const formChildrenArea = document.getElementById(
        newBlockModal.dataset.formTarget,
      );

      const id = newBlockModal.dataset.id;
      newBlockModal.dataset.id = id - 0 + 1;

      const parentId = blockParent.dataset.id;

      const blockChildArea = blockParent.querySelector(
        blockParent.dataset.childrenSelector,
      );

      const req = new Request(
        `${urls.blockCreate}?id=${id}&name=${newBlockName.value}&parentid=${parentId}`,
        {
          method: "GET",
        },
      );
      fetch(req)
        .then((resp) => resp.text())
        .then((html) => {
          log(html);
          newBlockWkSpace.innerHTML = html;
          const template = newBlockWkSpace.getElementsByTagName("template")[0];
          const [block, form] = template.content.children;
          log(block, form);
          formChildrenArea.appendChild(form);
          blockChildArea.appendChild(block);

          formEvents(form);
          newBlockWkSpace.innerHTML = "";

          newBlockModal.style.display = "none";
        })
        .catch(console.error);
    });
  };

  const pageEvents = () => {
    const doSave = (data, id) => {
      const req = new Request(`${urls.save}?id=${id}`, {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
          "Content-Type": "application/json",
        },
      });
      fetch(req)
        .then((resp) => {
          log(resp);
        })
        .catch(console.error);
    };
    const saveButton = document.getElementById("editor~save");
    saveButton.addEventListener("click", (_) => {
      const form = document.querySelector(".block-form");
      const data = getFormData(form, true);

      const id = saveButton.dataset.id;

      doSave(data, id);
    });
  };

  formEvents(document);
  modalEvents();
  pageEvents();
});
