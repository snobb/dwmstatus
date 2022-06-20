/* Source: dwm.suckless.org/dwmstatus/getvol.c*/
#include <alsa/asoundlib.h>
#include <alsa/control.h>

#define MUTED 0

int get_volume(void)
{
    long min, max, volume = 0;
    int status = !MUTED;
    snd_mixer_t *handle;
    snd_mixer_selem_id_t *sid;
    const char *card = "default";
    const char *selem_name = "Master";

    snd_mixer_open(&handle, 0);
    snd_mixer_attach(handle, card);
    snd_mixer_selem_register(handle, NULL, NULL);
    snd_mixer_load(handle);

    snd_mixer_selem_id_alloca(&sid);
    snd_mixer_selem_id_set_index(sid, 0);
    snd_mixer_selem_id_set_name(sid, selem_name);
    snd_mixer_elem_t* elem = snd_mixer_find_selem(handle, sid);

    snd_mixer_selem_get_playback_switch(elem, SND_MIXER_SCHN_MONO, &status);
    if (status == MUTED) {
        return -1;
    }

    snd_mixer_selem_get_playback_volume_range(elem, &min, &max);
    snd_mixer_selem_get_playback_volume(elem, 0, &volume);
    snd_mixer_close(handle);

    return ((double)volume / max) * 100;
}
