function neo_brutal_dropdown_toggle(id) {
	const menu = document.getElementById(id);
	if (menu == null) {
		console.error(`neo_brutal_dropdown_toggle: element ${id} not found`);
		return;
	}
	menu.classList.toggle("hidden");
}
