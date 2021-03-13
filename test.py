# Create a CommandSet for your remote control
# GPIO for the IR receiver: 23
# GPIO for the IR transmitter: 22
from ircodec.command import CommandSet
controller = CommandSet(emitter_gpio=18, receiver_gpio=17, description='Eufy')

# Add the volume up key
controller.add('start')
# Connected to pigpio
# Detecting IR command...
# Received.

# Send the volume up command
controller.emit('start')

# # Remove the volume up command
# controller.remove('volume_up')
#
# # Examine the contents of the CommandSet
# controller
# # CommandSet(emitter=22, receiver=23, description="TV remote")
# # {}
#
# # Save to JSON
# controller.save_as('tv.json')
#
# # Load from JSON
# new_controller = CommandSet.load('another_tv.json')
